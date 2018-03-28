HPC Object Tagger
=================

Golang lambda function to apply tags to newly created s3 objects based upon
suffix.

Specifically, the following mappings are maintained:

`*.fastq.gz` => `filetype: fastq`
`*.bam` => `filetype: bam`

If the file has a size > 50MB, it'll get `filetype: largefile` applied.

Cost Estimate
-------------

Assuming this takes 900ms for hit, then it's $2.07 / 1MM (128MB-s / 1024 * 0.9 * 0.00001667 * 1000000  + 0.20/1MM requests)
Assuming this takes 200ms for miss, then it's $0.62 / 1MM (128MB-s / 1024 * 0.2 * 0.00001667 * 1000000  + 0.20/1MM requests)

At 2.2MM requests / day (where 90% are misses) this should be <= cost for GoCD
instances to manage this.
