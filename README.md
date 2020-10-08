# Minecraft 1.15.2 Server

This is a minecraft server written from scratch in Go!

The goal is to create a full SMP vanilla server which will hopefully not be as heavy as the official one...

If you have any questions feel free to message me on discord `itay#2805`

## Features

Currently the focus is getting creative mode only fully working with all the block updates and such, 
afterwards we can get into implementing proper gameplay features that are needed for survival.

Once that is going we can get into implementing a proper anti-cheat into the game and so on, but that is 
far far away :)

## Basics
* Chunk sync
    * Chunk download
    * Chunk updates (block placing/breaking) 
    * Proper light calculations
* Entity & Players sync
    * Proper player spawning with Player info
    * All movement (including crouching and sprinting)
    * Equipment (sends a fake item that looks the same to not reveal the real properties of the item)
    * Animation (hand swing)

## Special blocks
* Proper placement of redstone and normal torch
