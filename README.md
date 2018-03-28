HPC Object Tagger
=================

Golang lambda function to apply tags to newly created s3 objects based upon
suffix.

Specifically, if a key ends in `.fastq.gz`,
it'll be tagged as `filetype: fastq`, if it ends in `.bam`, it'll be tagged as
`filetype: bam`. Lifecycle policies then transition these objects to Glacier.
