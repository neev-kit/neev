# Security Policy

## Reporting Security Vulnerabilities

**Please do NOT open a public GitHub issue for security vulnerabilities.**

If you discover a security vulnerability in Neev, please report it responsibly by emailing the maintainers or using GitHub's private vulnerability reporting feature.

### How to Report

1. **GitHub Security Advisory** (Recommended)
   - Go to [Security → Report a vulnerability](https://github.com/neev-kit/neev/security/advisories)
   - This creates a private vulnerability report visible only to maintainers
   - You'll be notified when the issue is resolved

2. **Direct Contact**
   - Email: Report via GitHub (see link above)
   - Contact: maintainers through private message on GitHub

### What to Include

Please provide:
- Description of the vulnerability
- Affected version(s)
- Steps to reproduce (if applicable)
- Potential impact
- Any known mitigations

## Response Timeline

We aim to:
- **Acknowledge** receipt within 24 hours
- **Assess** severity and impact within 48 hours
- **Provide** initial timeline for fix within 5 business days
- **Publish** security advisory after fix is released

## Security Maintenance

### Supported Versions

| Version | Supported |
|---------|-----------|
| 1.x     | ✅ Current |
| 0.x     | ⚠️ Legacy |

We provide security patches for the current major version. Legacy versions receive critical security fixes only.

### Dependencies

Neev has minimal dependencies to reduce attack surface:
- Reduces dependency vulnerabilities
- Easier security audits
- Faster updates when needed

We monitor dependencies via:
- GitHub Dependabot alerts
- Regular security audits
- Community reports

### Binary Integrity

All released binaries are:
- Built from tagged releases in CI/CD
- Signed checksums provided
- Verifiable builds

To verify binary integrity:
```bash
# Download checksum file
curl -O https://github.com/neev-kit/neev/releases/download/v1.0.0/checksums.txt

# Verify binary
sha256sum -c checksums.txt
```

## Security Best Practices for Users

### Installation
- ✅ Use official releases from [GitHub Releases](https://github.com/neev-kit/neev/releases)
- ✅ Verify checksums before running
- ✅ Use latest version for security patches
- ❌ Don't run arbitrary code in `.neev/` directories

### Configuration
- ✅ Keep Neev updated
- ✅ Review `.neev/` files before committing to git
- ✅ Use strong permissions on `.neev/foundation/` files
- ❌ Don't store secrets in blueprint files

### Data Protection
- `.neev/` files may contain architecture details
- Don't commit sensitive credentials
- Review public/private access controls in `neev.yaml`

## Coordinated Disclosure

We follow responsible disclosure practices:
1. Vulnerability reported to us privately
2. We assess and develop a fix
3. Security patch released
4. Public advisory published with credit to reporter (if desired)

## Security Updates

Stay informed about security updates:
- Watch releases: [GitHub Releases](https://github.com/neev-kit/neev/releases)
- Subscribe: GitHub repository notifications
- Follow: [@neev-kit](https://github.com/neev-kit) on GitHub

## Public Advisories

Security advisories are published at:
- [GitHub Security Advisories](https://github.com/neev-kit/neev/security/advisories)
- [GitHub Releases](https://github.com/neev-kit/neev/releases) (changelog)

## Questions?

- **Security concerns**: Use private vulnerability report (GitHub)
- **Policy questions**: Open an issue or discussion
- **General questions**: See [FAQ.md](FAQ.md)

---

Thank you for helping keep Neev secure! We appreciate your responsible disclosure.
