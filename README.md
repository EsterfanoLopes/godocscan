# godocscan
Read documents and Go

## Instalation On Ubuntu

1. Install Tesseract
``` bash
apt install tesseract-ocr
apt install libtesseract-dev
```

2. Install required language
``` bash
apt-get install tesseract-ocr-por
```

3. Install go mod
``` bash
go mod vendor
```

4. Run
``` bash
make run
```