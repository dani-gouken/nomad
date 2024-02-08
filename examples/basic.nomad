fun equation(int x, int y, int z): void {
    float delta: x*x - 4 * a * c;
    float s1: (-y - sqrt(delta)) / 2 *a;
    float s2: (-y + sqrt(delta)) / 2 *a;

    print_ln("The solution 1 is: "+ s1);
    print_ln("The solution 2 is: "+ s2);
}

func get_message() string {
    int age: 0;
    string message: "";

    if age < 18 {
        message: "Vous etes mineur";
    } else if age < 50 {
        message: "Vous etes adulte";
    }
    return message;  
}

type Gender: enum{male, female};

type Person: struct {
    pub {
        nom: string;
        prenon: string;
        gender: gender = gender.male;
    }

    pub func set_name(self, name) string {
        self.nom = name;
        return name;
    }
    
    func set_gender(Gender gender) void {
        self.gender = gender;
    }

    func get_name() string {
        return name;
    }
}
func Person_new(
    string name, string preson
) Person {
    return Person {
        name: name,
        prenom: person
    }
}

Person person = {
    nom: "",
    prenom: ""
};

print_ln person.set_name(lol);
print_ln person.get_name(lol);
