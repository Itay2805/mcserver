from ctypes import *
from datetime import datetime
import os.path as path
import urllib.request
import json

# TODO: switch to using go:generate

# This has all of the data we need
import stringcase

with urllib.request.urlopen("https://raw.githubusercontent.com/PrismarineJS/minecraft-data/master/data/pc/1.15.2/items.json") as f:
    items = json.load(f)

print("// Code generated by scripts/generate_items.go; DO NOT EDIT.")
print("// This file was generated by robots at")
print("// " + str(datetime.now()))
print()

print("package item")
print()

# Item variables
for item in items:
    name = stringcase.pascalcase(item['name'])
    print(f"var {name} = &Item{{")
    print(f"\tID: {item['id']},")
    print(f"\tName: \"{item['name']}\",")
    print(f"\tStackSize: {item['stackSize']},")
    print("}")
print()

# Item by id lookup
print("var items = [...]*Item{")
for item in items:
    name = stringcase.pascalcase(item['name'])
    print(f"\t{item['id']}: {name},")
print("}")
print()
