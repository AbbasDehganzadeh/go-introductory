## How to install
first, if you don't have golang, install that via this site.
[Go Site](https://go.dev)

then clone the repo
``` bash
  git clone https://github.com/AbbasDehganzadeh/go-introductory.git
```

,and go-to this directory.
``` bash
  cd go-introductory/mygit_label
```

## How to run
run the following command;
``` bash
  go run . --items=issues,prs --patterns=^a*,b*z$ --labels=lbl,llb --topics=top-1,top-2 output=data.json --normal
```
NOTE: are parameters are optional..
NOTE: parameters are saved for later usage.

### parameters
--items: whether you wanna see issues, prs, or both

--patterns: it's regex, and matches in titles, bodies, and comments
--labels: shows issues, and prs which have these labels
--topics: searches through repos which have these topics

~ if three of them are empty, it shows alll.
~ in labels, and topics numbers, and characters are removed.

--output: the file name for result data, if empty doesn't save;
~ format is JSON;
  format: prints the data in three format, --short, --normal, --verbose
~ One of them must be passed

###### END
