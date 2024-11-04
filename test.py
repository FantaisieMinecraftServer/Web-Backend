import json

import requests

# APIエンドポイントのURL
url = "http://localhost:8080/v2/items/"

# サンプルデータ
items = [
    {
        "id": "wooden_dagger",
        "type": "weapon",
        "name": "木の短剣",
        "lore": ["説明がないよ"],
        "rarity": 1,
        "max_stack_size": 1,
        "item_id": "stick",
        "custom_model_data": 10000,
        "prices": {"purchase": 0, "selling": 0, "can_selling": False},
        "data": {"group": "dagger"},
    },
    {
        "id": "wooden_sword",
        "type": "weapon",
        "name": "木の刀剣",
        "lore": ["説明がないよ"],
        "rarity": 1,
        "max_stack_size": 1,
        "item_id": "stick",
        "custom_model_data": 20000,
        "prices": {"purchase": 0, "selling": 0, "can_selling": False},
        "data": {"group": "sword"},
    },
    {
        "id": "wooden_spear",
        "type": "weapon",
        "name": "木の槍",
        "lore": ["説明がないよ"],
        "rarity": 1,
        "max_stack_size": 1,
        "item_id": "stick",
        "custom_model_data": 30000,
        "prices": {"purchase": 0, "selling": 0, "can_selling": False},
        "data": {"group": "spear"},
    },
    {
        "id": "wooden_hammer",
        "type": "weapon",
        "name": "ウッドハンマー",
        "lore": ["説明がないよ"],
        "rarity": 1,
        "max_stack_size": 1,
        "item_id": "stick",
        "custom_model_data": 40000,
        "prices": {"purchase": 0, "selling": 0, "can_selling": False},
        "data": {"group": "hammer"},
    },
    {
        "id": "wooden_wand",
        "type": "weapon",
        "name": "木の杖",
        "lore": ["説明がないよ"],
        "rarity": 1,
        "max_stack_size": 1,
        "item_id": "stick",
        "custom_model_data": 50000,
        "prices": {"purchase": 0, "selling": 0, "can_selling": False},
        "data": {"group": "wand"},
    },
    {
        "id": "wooden_bow",
        "type": "weapon",
        "name": "木の弓",
        "lore": ["一本の枝から作られた槍。", "しなやかで扱いやすい。"],
        "rarity": 1,
        "max_stack_size": 1,
        "item_id": "stick",
        "custom_model_data": 60000,
        "prices": {"purchase": 0, "selling": 0, "can_selling": False},
        "data": {"group": "bow"},
    },
]

# POSTリクエストを送信
for item in items:
    response = requests.post(
        url,
        json=item,
    )

    # レスポンスの表示
    if response.status_code == 201:
        print(
            "アイテム作成成功:",
            json.dumps(response.json(), ensure_ascii=False, indent=2),
        )
    else:
        print("アイテム作成失敗:", response.status_code, response.text)
