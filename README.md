# Minecraft 1.15.2 Server

This is a minecraft server written from scratch in Go!

The goal is to create a full SMP vanilla server which will hopefully not be as heavy as the official one...

## Features
* World
    * Chunk sync
        * currently, no backstore/real generation
        * Block breaking
        * Lighting
    * Entity sync
        * All entity movement
        * Proper player spawning with Player info
