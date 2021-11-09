package config

type Config struct {
	Targets                    []string `config:"targets,required" yaml:"targets"`
	FirestoreProjectOnEmulator string   `config:"firestore_project_on_emulator,required" yaml:"firestore_project_on_emulator"`
	FirestoreEmulatorHost      string   `config:"firestore_emulator_host,required" yaml:"firestore_emulator_host"`
}
