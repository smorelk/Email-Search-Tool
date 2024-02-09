#!/bin/bash
#
# Download and extract the email directory database for index later.
#
echo "Downloading archive"
wget http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz
echo "Done!"

echo "Extracting archive..."
tar -xf ./enron_mail_20110402.tgz
echo "Done!"


if [[ ! -d ./data ]]; then
    echo "Creating data directory for Zincsearch..."
    mkdir data
    echo "Applying proper permissions to ./data"
    chmod a+rwx ./data
fi

echo "Run the indexer now."