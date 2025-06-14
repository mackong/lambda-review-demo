import json


def lambda_handler(event, context):
    if "body" not in event:
        return "no body found"

    # 1. Parse GitHub Webhook
    payload = json.loads(event["body"])
    print(payload)
    return "ok"
