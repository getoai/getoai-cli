package installer

import "testing"

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		// EXE files
		{
			name:     "Windows EXE with setup",
			filename: "Cherry-Studio-1.7.13-x64-setup.exe",
			want:     "exe",
		},
		{
			name:     "Windows EXE simple",
			filename: "installer.exe",
			want:     "exe",
		},
		{
			name:     "Windows EXE uppercase",
			filename: "INSTALLER.EXE",
			want:     "exe",
		},

		// MSI files
		{
			name:     "Windows MSI",
			filename: "installer.msi",
			want:     "msi",
		},
		{
			name:     "Windows MSI uppercase",
			filename: "INSTALLER.MSI",
			want:     "msi",
		},

		// DMG files
		{
			name:     "macOS DMG arm64",
			filename: "Cherry-Studio-1.7.13-arm64.dmg",
			want:     "dmg",
		},
		{
			name:     "macOS DMG universal",
			filename: "Chatbox-1.18.3-universal.dmg",
			want:     "dmg",
		},

		// PKG files
		{
			name:     "macOS PKG",
			filename: "installer.pkg",
			want:     "pkg",
		},

		// DEB files
		{
			name:     "Linux DEB",
			filename: "Cherry-Studio_1.7.13_amd64.deb",
			want:     "deb",
		},
		{
			name:     "Linux DEB simple",
			filename: "package.deb",
			want:     "deb",
		},

		// AppImage files
		{
			name:     "Linux AppImage",
			filename: "jan-linux-x86_64-0.5.7.AppImage",
			want:     "appimage",
		},
		{
			name:     "Linux AppImage lowercase",
			filename: "app.appimage",
			want:     "appimage",
		},

		// Unknown/unsupported files
		{
			name:     "ZIP file",
			filename: "archive.zip",
			want:     "",
		},
		{
			name:     "TAR.GZ file",
			filename: "archive.tar.gz",
			want:     "",
		},
		{
			name:     "No extension",
			filename: "installer",
			want:     "",
		},
		{
			name:     "Empty string",
			filename: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectFileType(tt.filename)
			if got != tt.want {
				t.Errorf("detectFileType(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestGetFileNameFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Simple filename",
			url:  "https://example.com/installer.exe",
			want: "installer.exe",
		},
		{
			name: "Complex path",
			url:  "https://github.com/user/repo/releases/download/v1.0.0/app-1.0.0-x64.exe",
			want: "app-1.0.0-x64.exe",
		},
		{
			name: "With query parameters",
			url:  "https://download.example.com/app.dmg?version=1.0.0",
			want: "app.dmg",
		},
		{
			name: "Trailing slash",
			url:  "https://example.com/downloads/",
			want: "download",
		},
		{
			name: "Only domain",
			url:  "https://example.com",
			want: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFileNameFromURL(tt.url)
			if got != tt.want {
				t.Errorf("getFileNameFromURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}
