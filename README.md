# summary-generator-api

This is an API to summarize a given Japanese document.

Code for https://summary-generator.appspot.com/

## Algorithm

https://github.com/ramenjuniti/lexrank-mmr

## Usage

### Request

```
POST https://summary-generator.appspot.com/

# Request form-data
# {
#   "text": {input text},
#   "maxLines": {input maxLines (default 0)},
#   "maxCharacters": {input maxCharacters (default 0)},
#   "threshold": {input threshold (default 0.001)},
#   "tolerance": {input tolerance (default 0.0001)},
#   "damping": {input damping (default 0.85)},
#   "lambda": {input lambda (default 1.0)}
# }
```

### Response

```
{
  "LineLimitedSummary": [], # Enter here the line limited summary data.
  "CharacterLimitedSummary": [] # Enter here the character limited summary data.
}
```

## LICENSE

This sotfware is released under the MIT License, see LICENSE
