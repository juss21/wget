#!/bin/bash

echo 'Try to run the following command "./wget https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg"'
echo ""
go run . https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
echo 'Did the program download the file "EMtmPFLWkAA8CIS.jpg"?'

echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command with a link at your choice "./wget <https://link_of_your_choice.com>"'
echo ""
go run . https://i.pinimg.com/originals/e2/c6/1c/e2c61c6a0ca0f02097751bb19cae5f50.jpg
echo 'Did the program download the expected file?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command "./wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz"'
go run . https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
echo
echo 'Did the program download the file "go1.16.3.linux-amd64.tar.gz"?'
echo 'Did the program displayed the start time?'
echo 'Did the start time and the end time respected the format? (yyyy-mm-dd hh:mm:ss)'
echo 'Did the program displayed the status of the response? (200 OK)'
echo 'Did the Program displayed the content length of the download?'
echo 'Is the content length displayed as raw (bytes) and rounded (Mb or Gb)?'
echo 'Did the program displayed the name and path of the file that was saved?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to download a big file, for example: "./wget http://ipv4.download.thinkbroadband.com/100MB.zip"'
go run . http://ipv4.download.thinkbroadband.com/100MB.zip
echo
echo 'Did the program download the expected file?'
echo 'While downloading, did the progress bar show the amount that is being downloaded? (KiB or MiB)'
echo 'While downloading, did the progress bar show the percentage that is being downloaded?'
echo 'While downloading, did the progress bar show the time that remains to finish the download?'
echo 'While downloading, did the progress bar progressed smoothly (kept up with the time that the download took to finish)?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command, "./wget -O=test_20MB.zip http://ipv4.download.thinkbroadband.com/20MB.zip"'
go run . -O=test_20MB.zip http://ipv4.download.thinkbroadband.com/20MB.zip
echo 'Did the program downloaded the file with the name "test_20MB.zip"?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command, "./wget -O=test_20MB.zip -P=~/Downloads/ http://ipv4.download.thinkbroadband.com/20MB.zip"'
go run . -O=test_20MB.zip -P=~/Downloads/ http://ipv4.download.thinkbroadband.com/20MB.zip
echo 'Can you see the expected file in the "~/Downloads/" folder?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command, "./wget --rate-limit=300k http://ipv4.download.thinkbroadband.com/20MB.zip"'
go run . --rate-limit=300k http://ipv4.download.thinkbroadband.com/20MB.zip
echo 'Was the download speed always lower than 300KB/s?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command, "./wget --rate-limit=700k http://ipv4.download.thinkbroadband.com/20MB.zip"'
go run . --rate-limit=700k http://ipv4.download.thinkbroadband.com/20MB.zip
echo 'Was the download speed always lower than 700KB/s?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command, "./wget --rate-limit=2M http://ipv4.download.thinkbroadband.com/20MB.zip"'
go run . --rate-limit=2M http://ipv4.download.thinkbroadband.com/20MB.zip
echo 'Was the download speed always lower than 2MB/s?'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to create a text file with the name "downloads.txt" and save into it the links below. Then run the command "./wget -i=downloads.txt"'
go run . -i=downloads.txt
echo
echo 'Did the program download all the files from the downloads.txt file? (EMtmPFLWkAA8CIS.jpg, 20MB.zip, 10MB.zip)'
echo 'Did the downloads occurred in an asynchronous way? (tip: look to the download order)'
echo "Press enter to continue"
read "Press enter to continue"

echo 'Try to run the following command, "./wget -B http://ipv4.download.thinkbroadband.com/20MB.zip"'
go run . -B http://ipv4.download.thinkbroadband.com/20MB.zip
echo
echo 'Output will be written to ‘wget-log’.'
echo 'Did the program output the statement above?'
echo 'Was the download made in "silence" (without displaying anything to the terminal)?'
echo "Press enter to continue"
read "Press enter to continue"

echo Try to run the following command "./wget --mirror http://corndog.io/", then try to open the "index.html" with a browser
go run . --mirror http://corndog.io/
echo Is the site working?
echo "Press enter to continue"
read "Press enter to continue"

echo Try to run the following command "./wget --mirror https://oct82.com/", then try to open the "index.html" with a browser
go run . --mirror https://oct82.com/
echo Is the site working?
echo "Press enter to continue"
read "Press enter to continue"

echo Try to run the following command "./wget --mirror --reject=gif https://oct82.com/", then try to open the "index.html" with a browser
go run . --mirror --reject=gif https://oct82.com/
echo Did the program download the site without the GIFs?
echo "Press enter to continue"
read "Press enter to continue"

echo Try to run the following command "./wget --mirror https://trypap.com/", then use the command "ls" to see the file system of the created folder.
go run . --mirror https://trypap.com/
echo css  img  index.html
echo Does the created folder has the same fs as above?
echo "Press enter to continue"
read "Press enter to continue"

echo Try to run the following command "./wget --mirror -X=/img https://trypap.com/", then use the command "ls" to see the file system of the created folder.
go run . --mirror -X=/img https://trypap.com/
echo css index.html
echo Does the created folder has the same fs as above?
echo "Press enter to continue"
read "Press enter to continue"

echo Try to run the following command "./wget --mirror https://theuselessweb.com/"
go run . --mirror https://theuselessweb.com/
echo Is the site working?
echo "Press enter to continue"
read "Press enter to continue"

echo Try to mirror a website of your own choice "./wget --mirror https://www.postimees.ee"
go run . --mirror https://www.postimees.ee
echo Did the program mirror the website?
echo
echo 'END OF SCRIPT, TO EXIT PLEASE PRESS ENTER'
echo "Press enter to continue"
read "Press enter to continue"
