fun equation(int x, int y, int z): void {
    float delta = x*x - 4 * a * c;

    float s1 = (-y - math.sqrt(delta)) / 2 *a;
    float s2 = (-y + math.sqrt(delta)) / 2 *a;

    print_ln("The solution 1 is: "+ s1);
    print_ln("The solution 2 is: "+ s2);
}

type person = {
    age: int,
    name: string
};

impl stringable on person {
    func to_string(self): string {
        return self.name   
    }
}
func main {
    string name = io::std::readline_str();
    person p = person { age: 10, name: "daniel" };

    print_ln(person);
}