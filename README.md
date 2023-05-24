## detectdups
Here is an quick attempt to build a command line tool to detect duplicate files in given directory. When 2 or more files have exact same content, the tool considers them as duplicate and reports them.

### Sample run
```
$ detectdups data/
INFO[0000] processing dir: data/
INFO[0000] running detectdups
processing done
--final dupgs--
[data/orig data/orig-copy]
[data/orig-endcharchange data/orig-endcharchange-copy]
[data/orig-plus1 data/orig-plus1-copy]
--final dupgs-- #files with dups: 3
--final dupgs-- dupBytes: 304
INFO[0000] detectdups exiting
```
In above output all the lines with files listed in `[]` are duplicate of each other. For example `data/orig` and `data/orig-copy` are duplicate of each other.
