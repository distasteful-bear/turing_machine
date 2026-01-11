
# Deploy to Cloud Ru
gcloud run deploy turing-machine \
  --platform managed \
  --region us-east5 \
  --source . \
  --allow-unauthenticated
