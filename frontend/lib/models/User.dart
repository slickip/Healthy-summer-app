class User {
  final String email;
  final String displayName;

  User({required this.email, required this.displayName});

  factory User.fromJson(Map<String, dynamic> json) {
    return User(email: json['email'], displayName: json['display_name']);
  }
}
