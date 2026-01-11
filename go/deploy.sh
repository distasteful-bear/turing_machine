
# Deploy to Cloud Ru
gcloud run deploy turing-machine \
  --platform managed \
  --region us-east4 \
  --source . \
  --allow-unauthenticated \
  --min-instances 0 \
  --max-instances 1
