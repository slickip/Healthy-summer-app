import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class ChallengeDetailScreen extends StatefulWidget {
  const ChallengeDetailScreen({super.key});

  @override
  State<ChallengeDetailScreen> createState() => _ChallengeDetailScreenState();
}

class _ChallengeDetailScreenState extends State<ChallengeDetailScreen> {
  final ApiService api = ApiService();
  List<dynamic> leaderboard = [];
  Map<String, dynamic>? challenge;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    challenge =
        ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    fetchLeaderboard();
  }

  Future<void> fetchLeaderboard() async {
    final id = challenge?['id'];
    if (id == null) return;
    final result = await api.getChallengeLeaderboard(id);
    if (result != null) {
      setState(() {
        leaderboard = result;
      });
    }
  }

  Future<void> joinChallenge() async {
    final success = await api.joinChallenge(challenge!['id']);
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(success ? 'Joined!' : 'Failed to join')),
    );
    if (success) fetchLeaderboard();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(challenge?['title'] ?? 'Challenge Detail')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              challenge?['description'] ?? '',
              style: const TextStyle(fontSize: 16),
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: joinChallenge,
              child: const Text('Join Challenge'),
            ),
            const SizedBox(height: 16),
            const Text(
              'Leaderboard:',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            Expanded(
              child: ListView.builder(
                itemCount: leaderboard.length,
                itemBuilder: (context, index) {
                  final entry = leaderboard[index];
                  return ListTile(
                    title: Text('User ID: ${entry['user_id']}'),
                    trailing: Text('${entry['progress']} pts'),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
