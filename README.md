# Fine tune gpt

Set your `OPENAI_API_KEY` environment variable by adding the following line
into your shell initialization script (e.g. .bashrc, zshrc, etc.) or running
it in the command line before the fine-tuning command:

```bash
export OPENAI_API_KEY="<OPENAI_API_KEY>"
```

## Install

To build and to install the app you can use the `Makefile`. For help you can
run

```bash
make help
```

To install the app you can run

```bash
make compile && make install
```

## Usage

To generate synthetic data based on params you can run:

```bash
ftgpt generate-synt-data \
  --template=./example-data/synt/template.txt \
  --max-completions=60 \
  -- cogintives=./example-data/synt/cognitives.txt semantics=./example-data/synt/semantics.txt goals=./example-data/synt/goals.txt
```

To generate json file from producerd syntehtic data you can run:

```bash
ftgpt generate-json -p ./data/prompts/ -c ./data/completions/
```

Prepare data will generate jsonl file that you can use to create the fine tune:

```
openai tools fine_tunes.prepare_data -f <LOCAL_FILE>
```

To create the fine tune you can run:

```
openai api fine_tunes.create -t <TRAIN_FILE_ID_OR_PATH> -m <BASE_MODEL>
```
