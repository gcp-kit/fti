package config

type Config struct {
	Targets                    []string `config:"targets,required" yaml:"targets"`
	FirestoreProjectOnEmulator string   `config:"firestore_project_on_emulator" yaml:"firestore_project_on_emulator"`
	FirestoreEmulatorHost      string   `config:"firestore_emulator_host" yaml:"firestore_emulator_host"`
}
