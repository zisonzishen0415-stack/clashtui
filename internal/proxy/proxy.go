package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"clashtui/internal/state"
)

// SetSystemProxy sets system proxy for GUI apps and creates shell script
func SetSystemProxy(port int) {
	portStr := fmt.Sprintf("%d", port)

	if hasGsettings() {
		exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "manual").Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "host", "127.0.0.1").Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "port", portStr).Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "host", "127.0.0.1").Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "port", portStr).Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "host", "127.0.0.1").Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "port", portStr).Run()
		exec.Command("gsettings", "set", "org.gnome.system.proxy", "ignore-hosts",
			"['localhost', '127.0.0.0/8', '::1']").Run()
	}

	if hasKwriteconfig() {
		kwriteCmd := getKwriteconfigCmd()
		exec.Command(kwriteCmd, "--file", "kioslaverc", "--group", "Proxy Settings",
			"--key", "ProxyType", "1").Run()
		exec.Command(kwriteCmd, "--file", "kioslaverc", "--group", "Proxy Settings",
			"--key", "httpProxy", fmt.Sprintf("http://127.0.0.1:%d", port)).Run()
		exec.Command(kwriteCmd, "--file", "kioslaverc", "--group", "Proxy Settings",
			"--key", "httpsProxy", fmt.Sprintf("http://127.0.0.1:%d", port)).Run()
		exec.Command(kwriteCmd, "--file", "kioslaverc", "--group", "Proxy Settings",
			"--key", "NoProxyFor", "localhost,127.0.0.1").Run()
	}

	createChromeProxyDesktop(port)

	createProxyScript(port, true)

	state.SaveState(state.NetworkState{
		Mode: state.ModeSystemProxy,
		Port: port,
	})
}

// UnsetSystemProxy clears system proxy
func UnsetSystemProxy() {
	if hasGsettings() {
		exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none").Run()
	}

	if hasKwriteconfig() {
		kwriteCmd := getKwriteconfigCmd()
		exec.Command(kwriteCmd, "--file", "kioslaverc", "--group", "Proxy Settings",
			"--key", "ProxyType", "0").Run()
	}

	removeChromeProxyDesktop()

	createProxyScript(0, false)

	state.SaveState(state.NetworkState{
		Mode: state.ModeOff,
	})
}

// createProxyScript creates ~/.config/clashtui/proxy.sh for user to source
// The script is smart - it checks if mihomo is running before setting env vars
func createProxyScript(port int, enable bool) {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "clashtui")
	os.MkdirAll(configDir, 0755)

	scriptPath := filepath.Join(configDir, "proxy.sh")

	var content string
	if enable {
		// Smart script that only sets proxy if mihomo is actually running
		content = fmt.Sprintf(`#!/bin/sh
# ClashTUI proxy env vars - run: source ~/.config/clashtui/proxy.sh
# Smart: only sets proxy if mihomo is running

MIHOMO_PORT=%d

# Check if mihomo is running by testing the API
if curl -s --connect-timeout 1 http://127.0.0.1:9090/version >/dev/null 2>&1; then
    export HTTP_PROXY=http://127.0.0.1:$MIHOMO_PORT
    export HTTPS_PROXY=http://127.0.0.1:$MIHOMO_PORT
    export ALL_PROXY=socks5://127.0.0.1:$MIHOMO_PORT
    export NO_PROXY=localhost,127.0.0.1
    echo "Proxy enabled: 127.0.0.1:$MIHOMO_PORT"
else
    unset HTTP_PROXY
    unset HTTPS_PROXY
    unset ALL_PROXY
    unset NO_PROXY
    unset no_proxy
    echo "Proxy disabled: mihomo not running"
fi
`, port)
	} else {
		content = `#!/bin/sh
# ClashTUI proxy cleared - run: source ~/.config/clashtui/proxy.sh
unset HTTP_PROXY
unset HTTPS_PROXY
unset ALL_PROXY
unset NO_PROXY
unset no_proxy
echo "Proxy disabled"
`
	}

	os.WriteFile(scriptPath, []byte(content), 0755)
}

func hasGsettings() bool {
	_, err := exec.LookPath("gsettings")
	return err == nil
}

func hasKwriteconfig() bool {
	_, err := exec.LookPath("kwriteconfig6")
	if err == nil {
		return true
	}
	_, err = exec.LookPath("kwriteconfig5")
	return err == nil
}

func getKwriteconfigCmd() string {
	if _, err := exec.LookPath("kwriteconfig6"); err == nil {
		return "kwriteconfig6"
	}
	return "kwriteconfig5"
}

// createChromeProxyDesktop creates a desktop file for Chrome with proxy flag
// Chrome on Wayland ignores gsettings, needs explicit --proxy-server flag
func createChromeProxyDesktop(port int) {
	home, _ := os.UserHomeDir()
	appDir := filepath.Join(home, ".local/share/applications")
	os.MkdirAll(appDir, 0755)

	desktopPath := filepath.Join(appDir, "chrome-proxy.desktop")

	content := fmt.Sprintf(`[Desktop Entry]
Name=Chrome (Proxy)
Exec=google-chrome --proxy-server="http://127.0.0.1:%d"
Type=Application
Icon=google-chrome
`, port)

	os.WriteFile(desktopPath, []byte(content), 0755)
}

// removeChromeProxyDesktop removes the Chrome proxy desktop file
func removeChromeProxyDesktop() {
	home, _ := os.UserHomeDir()
	desktopPath := filepath.Join(home, ".local/share/applications", "chrome-proxy.desktop")
	os.Remove(desktopPath)
}

func GetSystemProxyPort() int {
	if hasGsettings() {
		output, _ := exec.Command("gsettings", "get", "org.gnome.system.proxy.http", "port").Output()
		portStr := strings.TrimSpace(string(output))
		if portStr != "" && portStr != "0" {
			var port int
			fmt.Sscanf(portStr, "%d", &port)
			return port
		}
	}
	return 0
}