# macOS Notarization

The `go-binary-release.yml` reusable workflow supports automatic code signing and notarization of macOS binaries via GoReleaser. This prevents Gatekeeper from quarantining binaries distributed through Homebrew or direct download.

## Prerequisites

- An active Apple Developer Program membership ($99/year, enrolled as Organization or Individual)
- A **Developer ID Application** certificate
- An **App Store Connect API key** with Developer role

## Step 1: Create a Developer ID Application Certificate

1. Go to [developer.apple.com/account/certificates](https://developer.apple.com/account/certificates) → **+**
2. Under **Developer ID**, select **Developer ID Application** → choose **G2 Sub-CA**
3. Generate a CSR on your Mac: open **Keychain Access → Certificate Assistant → Request a Certificate from a Certificate Authority** → save to disk
4. Upload the `.certSigningRequest` file → download the resulting `.cer`
5. Double-click the `.cer` to install it into your **login** keychain

## Step 2: Export the Certificate as `.p12`

1. Open **Keychain Access** → select **login** keychain → **My Certificates**
2. Right-click **Developer ID Application: Your Org (TEAMID)** → **Export** → save as `.p12` with a strong password

Base64-encode it:
```bash
base64 -i ~/Desktop/certificate.p12 | pbcopy
```

## Step 3: Create an App Store Connect API Key

1. Go to [appstoreconnect.apple.com](https://appstoreconnect.apple.com) → **Users and Access → Integrations → App Store Connect API**
2. Click **+** → name it (e.g. `notarization`) → role: **Developer** → **Generate**
3. Note the **Issuer ID** (top of page) and **Key ID** (next to your key)
4. Download the `.p8` private key file (shown only once)

Base64-encode the key:
```bash
base64 -i ~/Downloads/AuthKey_XXXXXXXXXX.p8 | pbcopy
```

## Step 4: Add GitHub Secrets

Add these secrets to each repo that uses `go-binary-release.yml`:

| Secret | Value |
|---|---|
| `MACOS_CERTIFICATE` | base64-encoded `.p12` certificate |
| `MACOS_CERTIFICATE_PWD` | password set when exporting `.p12` |
| `NOTARIZATION_ISSUER_ID` | Issuer ID from App Store Connect |
| `NOTARIZATION_KEY_ID` | Key ID from App Store Connect |
| `NOTARIZATION_KEY` | base64-encoded `.p8` private key |

## Step 5: Configure the Calling Workflow

```yaml
goreleaser:
  uses: <org>/release-foundry/.github/workflows/go-binary-release.yml@main
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
    HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
    MACOS_CERTIFICATE: ${{ secrets.MACOS_CERTIFICATE }}
    MACOS_CERTIFICATE_PWD: ${{ secrets.MACOS_CERTIFICATE_PWD }}
    NOTARIZATION_ISSUER_ID: ${{ secrets.NOTARIZATION_ISSUER_ID }}
    NOTARIZATION_KEY_ID: ${{ secrets.NOTARIZATION_KEY_ID }}
    NOTARIZATION_KEY: ${{ secrets.NOTARIZATION_KEY }}
```

## Step 6: Configure GoReleaser

Add to your `.goreleaser.yml`:

```yaml
notarize:
  macos:
    - enabled: true
      sign:
        certificate: "{{ .Env.MACOS_CERTIFICATE }}"
        password: "{{ .Env.MACOS_CERTIFICATE_PWD }}"
      notarize:
        issuer_id: "{{ .Env.NOTARIZATION_ISSUER_ID }}"
        key_id: "{{ .Env.NOTARIZATION_KEY_ID }}"
        key: "{{ .Env.NOTARIZATION_KEY }}"

homebrew_casks:
  - ...
    disable_quarantine: true
```

> **Note:** `go-binary-release.yml` runs on `macos-latest` because `codesign` is macOS-only. This adds ~3-4 minutes to release builds.
