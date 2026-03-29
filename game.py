import random
import time
import os

def clear():
    os.system('clear')

class Player:
    def __init__(self, name):
        self.name = name
        self.hp = 100
        self.max_hp = 100
        self.attack = 15
        self.gold = 0
        self.xp = 0
        self.level = 1

    def is_alive(self):
        return self.hp > 0

class Enemy:
    def __init__(self, name, hp, attack, gold_drop):
        self.name = name
        self.hp = hp
        self.attack = attack
        self.gold_drop = gold_drop

def battle(player):
    enemies = [
        Enemy("Buggy Script", 30, 5, 20),
        Enemy("Malware Bot", 50, 10, 50),
        Enemy("Kernel Ghost", 80, 18, 100)
    ]
    enemy = random.choice(enemies)
    print(f"--- PERINGATAN! {enemy.name} Muncul! ---")
    
    while enemy.hp > 0 and player.is_alive():
        print(f"\n{player.name} HP: {player.hp} | {enemy.name} HP: {enemy.hp}")
        action = input("Aksi: (1) Serang (2) Lari: ")
        
        if action == "1":
            dmg = random.randint(player.attack - 5, player.attack + 5)
            enemy.hp -= dmg
            print(f"Kamu menyerang {enemy.name} sebesar {dmg} damage!")
            
            if enemy.hp > 0:
                e_dmg = random.randint(enemy.attack - 3, enemy.attack + 3)
                player.hp -= e_dmg
                print(f"{enemy.name} membalas {e_dmg} damage!")
        elif action == "2":
            if random.random() > 0.5:
                print("Berhasil kabur!")
                return
            else:
                print("Gagal kabur! Musuh menyerang!")
                player.hp -= enemy.attack
        
    if player.is_alive():
        print(f"\nSelamat! {enemy.name} hancur. Kamu mendapat {enemy.gold_drop} Gold.")
        player.gold += enemy.gold_drop
        player.xp += 50
        if player.xp >= 100:
            player.level += 1
            player.max_hp += 20
            player.hp = player.max_hp
            player.attack += 5
            player.xp = 0
            print("LEVEL UP! Stat meningkat.")
        input("\nTekan Enter untuk lanjut...")

def main():
    clear()
    print("=== WELCOME TO INDOS CYBER-DUNGEON ===")
    name = input("Masukkan Nama Karakter: ")
    p = Player(name)
    
    while p.is_alive():
        clear()
        print(f"User: {p.name} | Level: {p.level} | HP: {p.hp}/{p.max_hp} | Gold: {p.gold}")
        print("-" * 40)
        print("1. Jelajahi Terminal (Cari Musuh)")
        print("2. Rest (Pulihkan HP - 30 Gold)")
        print("3. Keluar")
        
        choice = input("\nPilih tindakan: ")
        
        if choice == "1":
            battle(p)
        elif choice == "2":
            if p.gold >= 30:
                p.gold -= 30
                p.hp = p.max_hp
                print("HP Pulih sepenuhnya!")
                time.sleep(1)
            else:
                print("Gold tidak cukup!")
                time.sleep(1)
        elif choice == "3":
            print("Shutdown system...")
            break
            
    if not p.is_alive():
        print("\nGAME OVER. System Corrupted.")

if __name__ == "__main__":
    main()
a