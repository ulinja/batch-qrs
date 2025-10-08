# Batch QRs

Batch QRs is a tiny program which batch-creates QR code images from a list of URL strings.

The main focus here is **performance**: batch QR encoding is CPU-parallelized and avoids disk I/O entirely, making it really fast.

### Usage

The URL strings to be turned into QR code images are supplied as arguments to the program:

```bash
./batch_qrs URL...
```

The program outputs one base64-encoded PNG image per each URL, separated by newlines.

#### Example

```bash
# Outputs three base64-encoded PNG images
./batch_qrs https://example.com/ https://acme.org/ https://wikipedia.org/
```

### Building

To build Batch QRs, you need Go >= 1.24.7.
Build by running:

```bash
go build batch_qrs.go
```
