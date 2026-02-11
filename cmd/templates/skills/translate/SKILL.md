---
name: translate
description: Translate text between languages using free APIs.
tags: [translate, language, cross-platform]
---
# Translate

Translate text between languages using free translation APIs. No API key required.

## Method 1: MyMemory API (Free, No Key)

Translate text:
```
exec: curl -s "https://api.mymemory.translated.net/get?q=TEXT_TO_TRANSLATE&langpair=SOURCE_LANG|TARGET_LANG" | jq -r '.responseData.translatedText'
```

Examples:

English to Chinese:
```
exec: curl -s "https://api.mymemory.translated.net/get?q=Hello%20World&langpair=en|zh" | jq -r '.responseData.translatedText'
```

Chinese to English:
```
exec: curl -s "https://api.mymemory.translated.net/get?q=你好世界&langpair=zh|en" | jq -r '.responseData.translatedText'
```

English to Japanese:
```
exec: curl -s "https://api.mymemory.translated.net/get?q=Good%20morning&langpair=en|ja" | jq -r '.responseData.translatedText'
```

With email (higher rate limit):
```
exec: curl -s "https://api.mymemory.translated.net/get?q=TEXT&langpair=en|zh&de=your@email.com" | jq -r '.responseData.translatedText'
```

## Method 2: LibreTranslate (Self-Hosted or Public Instances)

If a public instance is available:
```
exec: curl -s -X POST "https://libretranslate.com/translate" -H "Content-Type: application/json" -d '{"q":"TEXT_TO_TRANSLATE","source":"en","target":"zh"}' | jq -r '.translatedText'
```

Detect language:
```
exec: curl -s -X POST "https://libretranslate.com/detect" -H "Content-Type: application/json" -d '{"q":"TEXT"}' | jq '.[0]'
```

List supported languages:
```
exec: curl -s "https://libretranslate.com/languages" | jq '.[].code'
```

## Method 3: Lingva Translate (Google Translate Frontend, No Key)

```
exec: curl -s "https://lingva.ml/api/v1/SOURCE_LANG/TARGET_LANG/TEXT_URL_ENCODED" | jq -r '.translation'
```

Example:
```
exec: curl -s "https://lingva.ml/api/v1/en/zh/Hello%20World" | jq -r '.translation'
```

## Translate a File

Read file, translate, save:
```
exec: TEXT=$(cat /path/to/input.txt | head -500 | python3 -c "import sys,urllib.parse; print(urllib.parse.quote(sys.stdin.read().strip()))") && curl -s "https://api.mymemory.translated.net/get?q=$TEXT&langpair=en|zh" | jq -r '.responseData.translatedText' > /path/to/output.txt
```

## Common Language Codes

- `en` — English
- `zh` — Chinese (Simplified)
- `zh-TW` — Chinese (Traditional)
- `ja` — Japanese
- `ko` — Korean
- `es` — Spanish
- `fr` — French
- `de` — German
- `ru` — Russian
- `ar` — Arabic
- `pt` — Portuguese
- `it` — Italian

## Notes

- MyMemory API: 5000 chars/day free (anonymous), 50000/day with email. Best for short texts.
- LibreTranslate: some public instances may be rate-limited or down. Self-host for reliability.
- Lingva: acts as a frontend to Google Translate; availability varies.
- URL-encode text with spaces: replace spaces with `%20` or use `+`.
- For long texts, split into chunks under the character limit.
- All methods use `curl` and `jq` — works on macOS, Linux, Windows (WSL).
