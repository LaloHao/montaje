import argparse
import sys
from pathlib import Path

# faster-whisper no necesita PyTorch; usa CTranslate2
from faster_whisper import WhisperModel

SRT_TEMPLATE = """{index}\n{start} --> {end}\n{text}\n\n"""

def fmt_ts(seconds: float) -> str:
    if seconds < 0:
        seconds = 0
    ms = int(round((seconds - int(seconds)) * 1000))
    s = int(seconds) % 60
    m = int(seconds // 60) % 60
    h = int(seconds // 3600)
    return f"{h:02d}:{m:02d}:{s:02d},{ms:03d}"


def write_srt(segments, out_path: Path):
    with out_path.open('w', encoding='utf-8') as f:
        for i, seg in enumerate(segments, start=1):
            f.write(SRT_TEMPLATE.format(index=i, start=fmt_ts(seg.start), end=fmt_ts(seg.end), text=seg.text.strip()))


def run(input_path: Path, output_path: Path, model_size: str, device: str, language: str | None):
    # device: "auto" intenta cuda->metal->cpu
    compute_type = "float16" if device in ("cuda", "auto") else "int8_float16"
    model = WhisperModel(model_size, device=device if device != "auto" else "auto", compute_type=compute_type)
    segments, info = model.transcribe(str(input_path), vad_filter=True, language=language)
    write_srt(list(segments), output_path)


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--input", required=True)
    ap.add_argument("--output", required=True)
    ap.add_argument("--model", default="large-v3")
    ap.add_argument("--device", default="auto", help="cuda|metal|cpu|auto")
    ap.add_argument("--language", default=None)
    args = ap.parse_args()

    inp = Path(args.input)
    outp = Path(args.output)
    outp.parent.mkdir(parents=True, exist_ok=True)


    if not inp.exists():
        print(f"Input not found: {inp}", file=sys.stderr)
        sys.exit(2)


    run(inp, outp, args.model, args.device, args.language)
    print(f"Wrote SRT → {outp}")