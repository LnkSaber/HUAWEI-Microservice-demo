package hello.service;

public class Person {
    private Gender gender;
    private String name;


    public Gender getGender() { return gender; }
    public void setGender(Gender gender) { this.gender = gender; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }

    public Person(Gender greeting, String name) {
        this.gender = greeting;
        this.name = name;
    }

    public Person() {
    }
}
