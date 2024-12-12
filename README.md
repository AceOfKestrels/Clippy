# Clippy
The new Clippy, your friend and helper! Just copy any question then press F8 to ask Clippy!
- Supports both multiple-choice and free form questions, though the latter will take longer to answer.
- Response is copied to clipboard automatically.

<br>

Make sure to add your [Gemini API key](https://aistudio.google.com/app/apikey) to the config!

### Config Documentation:
- `apiKey`: Your api key (required)
- `model`: The AI model to use (Default: `gemini-1.5-flash-latest`)
  - Suggested models:
  - `gemini-1.5-flash-latest`: High rate limits, but pretty slow and not that smart
  - `gemini-1.5-pro-latest`: Much lower rate limits, same speed, but much smarter
  - `gemini-1.5-flash-8b-latest`: Same rate limits, much faster, even less smart
  - `gemini-2.0-flash-exp`: Experimental model, average rate limits and speed, not sure about smartness