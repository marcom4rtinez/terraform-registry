package model

const DataPath = "data/providers"

type GPGPublicKey struct {
	KeyID          string `json:"key_id"`
	AsciiArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"source"`
	SourceURL      string `json:"source_url"`
}

type Platform struct {
	OS                  string `json:"os"`
	Arch                string `json:"arch"`
	Filename            string `json:"filename"`
	DownloadURL         string `json:"download_url"`
	ShasumsURL          string `json:"shasums_url"`
	ShasumsSignatureURL string `json:"shasums_signature_url"`
	Shasum              string `json:"shasum"`
	SigningKeys         struct {
		GPGPublicKeys []GPGPublicKey `json:"gpg_public_keys"`
	} `json:"signing_keys"`
}

type Version struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

type ProviderData struct {
	Versions []Version `json:"versions"`
}
