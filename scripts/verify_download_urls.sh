#!/bin/bash
# Script to verify all download URLs in the registry

echo "Verifying Download URLs..."
echo "=========================="
echo ""

# Extract all URLs from registry.go
urls=$(grep -o 'https://[^"]*' /Users/wyp/git/getoai-cli/internal/tools/registry.go | sort -u)

total=0
success=0
failed=0
redirect=0

while IFS= read -r url; do
    total=$((total + 1))

    # Skip non-specific URLs (package pages, not direct downloads)
    if [[ "$url" == *"/releases"$ ]] || \
       [[ "$url" == "https://cursor.sh" ]] || \
       [[ "$url" == "https://lmstudio.ai" ]] || \
       [[ "$url" == "https://jan.ai" ]] || \
       [[ "$url" == "https://msty.app" ]] || \
       [[ "$url" == "https://chatboxai.app" ]] || \
       [[ "$url" == "https://tableplus.com" ]] || \
       [[ "$url" == *"/download"$ ]]; then
        echo "⊘ SKIP: $url (package page, not direct download)"
        total=$((total - 1))
        continue
    fi

    # Test URL with HEAD request
    http_code=$(curl -o /dev/null -s -w "%{http_code}" -L --max-time 10 "$url" 2>/dev/null)

    if [ "$http_code" == "200" ]; then
        echo "✓ OK  : $url"
        success=$((success + 1))
    elif [ "$http_code" == "302" ] || [ "$http_code" == "301" ]; then
        echo "→ REDIR: $url (HTTP $http_code)"
        redirect=$((redirect + 1))
        success=$((success + 1))  # Redirects are OK
    elif [ "$http_code" == "000" ]; then
        echo "? TIMEOUT: $url"
        failed=$((failed + 1))
    else
        echo "✗ FAIL: $url (HTTP $http_code)"
        failed=$((failed + 1))
    fi
done <<< "$urls"

echo ""
echo "=========================="
echo "Summary:"
echo "  Total tested: $total"
echo "  Success: $success"
echo "  Redirects: $redirect"
echo "  Failed: $failed"
echo ""

if [ $failed -gt 0 ]; then
    echo "⚠️  Some URLs failed verification. Please check them manually."
    exit 1
else
    echo "✅ All URLs verified successfully!"
    exit 0
fi
