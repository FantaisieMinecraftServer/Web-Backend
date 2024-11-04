import json

import requests

# APIエンドポイントのURL
url = "http://localhost:8080/v2/items/"

# サンプルデータ
item_data = {
    "id": "wooden_sword",
    "type": "weapon",
    "name": "木の刀剣",
    "lore": ["木の枝から作られた刀剣。", "とても軽く切れ味が良い。"],
    "rarity": 1,
    "max_stack_size": 1,
    "item_id": "carrot_on_a_stick",
    "custom_model_data": 101,
    "prices": {"purchase": 100, "selling": 20, "can_selling": True},
    "data": {"group": "sword"},
}

# POSTリクエストを送信
response = requests.post(
    url,
    json=item_data,
)

# レスポンスの表示
if response.status_code == 201:
    print(
        "アイテム作成成功:", json.dumps(response.json(), ensure_ascii=False, indent=2)
    )
else:
    print("アイテム作成失敗:", response.status_code, response.text)
