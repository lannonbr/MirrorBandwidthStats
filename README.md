# Mirror Bandwidth Stats

This is a project to generate data for Clarkson's FOSS Mirror based off of data outputted
from dstat.

# Examples

For example, if you want to view the aggregate of July 2017, run the following:

```
$ ./MirrorBandwidthStats pretty_month $CSVPath/mirror-Jul-*-2017*
```

The first argument is the format which the data will be outputted.
it can be any of the following:

- pretty_{hour,day,month}: Pretty print based on the duration inputted
- csv_{hour,day,month}: print to a csv entry in a human format for the numbers (Ex: 3.3TB)
- csv_{hour,day,month}_raw: print to csv with raw numbers (Ex: 3317146885253)