extends Node

var players = 0

# Called when the node enters the scene tree for the first time.
func _ready():
	get_node("GameArea").player_add.connect(on_player_add)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	pass
	
func on_player_add():
	players += 1
	$HUD.update_player_count(players);
