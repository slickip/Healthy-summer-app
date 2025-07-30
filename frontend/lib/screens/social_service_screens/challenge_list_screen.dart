import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class ChallengeListScreen extends StatefulWidget {
  const ChallengeListScreen({super.key});

  @override
  State<ChallengeListScreen> createState() => _ChallengeListScreenState();
}

class _ChallengeListScreenState extends State<ChallengeListScreen> {
  final ApiService api = ApiService();
  List<dynamic> challenges = [];

  @override
  void initState() {
    super.initState();
    fetchChallenges();
  }

  Future<void> fetchChallenges() async {
    final result = await api.getChallenges();
    if (result != null) {
      setState(() {
        challenges = result;
      });
    }
  }

  void goToDetail(Map<String, dynamic> challenge) {
    Navigator.pushNamed(context, '/challenge_detail', arguments: challenge);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Challenges')),
      floatingActionButton: FloatingActionButton(
        onPressed: () => Navigator.pushNamed(context, '/create_challenge'),
        child: const Icon(Icons.add),
      ),
      body: challenges.isEmpty
          ? const Center(child: Text('No challenges found.'))
          : ListView.builder(
              itemCount: challenges.length,
              itemBuilder: (context, index) {
                final challenge = challenges[index];
                return ListTile(
                  title: Text(challenge['title']),
                  subtitle: Text(challenge['description']),
                  trailing: const Icon(Icons.arrow_forward),
                  onTap: () => goToDetail(challenge),
                );
              },
            ),
    );
  }
}
