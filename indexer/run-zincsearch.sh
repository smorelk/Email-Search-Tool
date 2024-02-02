if [[ ! -d ./data ]]; then
    echo "Creating data directory for Zincsearch..."
    mkdir data
    echo "Applying proper permissions"
    chmod a+rwx ./data
fi

echo "Starting Zincsearch docker container"
sudo docker run --rm -v $(pwd)/data:/data -e ZINC_DATA_PATH="/data" -p 4080:4080 \
    -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 \
    --name zincsearch public.ecr.aws/zinclabs/zincsearch:latest
