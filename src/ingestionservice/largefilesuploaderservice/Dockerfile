# Dockerfile snippet for using `gsutil` inside container # gsutil is very useful for uploading large files
# thanks to 'resumable upload' feature. The golang/python clients are more difficult to use with the resumable media upload
# so that's why I'll concurrently spawn gsutil processes to upload big files.


# Install Google Cloud tools - Debian https://cloud.google.com/storage/docs/gsutil_install#deb
ENV CLOUD_SDK_REPO="cloud-sdk-stretch"
RUN echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | \
    tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    apt-get update && apt-get install -y google-cloud-sdk

# Setup Google Service Account
COPY service-account.json /etc/
ENV GOOGLE_APPLICATION_CREDENTIALS="/etc/service-account.json"

RUN gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}