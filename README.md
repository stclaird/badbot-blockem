# BadBot-Blockem

Blockem is a simple tool written in GoLang for slurping Internet IP address blacklists and then outputting one giant list for your use. For example, you might use the list for creating a IP blacklists on your network / application firewall. Amazon WAF is one such Firewall and there are others for sure.

## Installation 'n Compile

It is fairly easy to install, you just need a GoLang compiler installed and have the this repo checked out to your machine. 
From there you compile the source go code into an executable.

```bash
go build blockem.go
```
Which if all goes to plan will output a executable called blockem

## Usage
Make sure the executable is executable by setting the executable bit before executing.

```bash
chmod +x blockem
```

Run it
```bash
./blockem
```
This will check the 'built in' IP address blacklists that blockem knows about. These all are available publicly on the internet. It will run through them concurrently (which is sort of mandatory for GoLang), and then spit out a file called badbots-blockem-ip-list.out
If you want to output to a different, and perhaps more professional sounding file, you can specify one with the fileout flag, like so:

```bash
./blockem -fileout <some-file.txt>
```

In the likely event that you will want to provide your own comma seperated file containing a list of blacklist URLs. Then we have you covered, you just need to specify the blacklist_urls flag:
```bash
./blockem -blacklist_urls blacklist_urls_file_in.txt
```
There is an example file included caled blacklist_urls_file_in.txt as an example.

## Contributing
Pull requests are welcome I guess if you are really bored.

## License
[MIT](https://choosealicense.com/licenses/mit/)
