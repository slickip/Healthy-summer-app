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
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text('Challenges', style: TextStyle(color: Colors.white)),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => Navigator.pushNamed(context, '/create_challenge'),
        backgroundColor: Colors.orange[700],
        child: const Icon(Icons.add),
      ),
      body: challenges.isEmpty
          ? const Center(child: CircularProgressIndicator())
          : ListView.builder(
              itemCount: challenges.length,
              itemBuilder: (context, index) {
                final challenge = challenges[index];
                return Card(
                  margin: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 8,
                  ),
                  child: ListTile(
                    title: Text(challenge['title']),
                    subtitle: Text(challenge['description']),
                    trailing: const Icon(Icons.arrow_forward),
                    onTap: () => goToDetail(challenge),
                  ),
                );
              },
            ),
    );
  }
}
