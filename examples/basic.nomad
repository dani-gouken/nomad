Dog () {
    import println from :io:stdout

    string name;
    
    func New() {

    }

    func<Dog> bark = () {
        printn("ouaf")
        return ::
    }

    func<bool> changeName = (string new_name) {
        ::name = new_name;
    }
}

Dog dog = Dog:new()
Dog dog = Dog:new()

func<float> equation =  () {
    func print = :io:stdout:print_ln

    print("Entrez a: \n")
    float a = :float.from_string(readline())
    print("Entrez b: \n")
    float b = :float.from_string(readline())
    print("Entrez c: \n")
    float c = :int.from_string(readline())

    int delta = :math.pow(b, 2) - (4*a*c)
    if res < 0 {
        print("Aucun résultat");
    }

}
equation()