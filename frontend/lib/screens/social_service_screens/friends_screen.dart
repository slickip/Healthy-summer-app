import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class FriendsScreen extends StatefulWidget {
  const FriendsScreen({super.key});

  @override
  State<FriendsScreen> createState() => _FriendsScreenState();
}

class _FriendsScreenState extends State<FriendsScreen> {
  final ApiService api = ApiService();
  List<dynamic> feed = [];

  @override
  void initState() {
    super.initState();
    fetchFeed();
  }

  Future<void> fetchFeed() async {
    final result = await api.getFriendsFeed();
    if (result != null) {
      setState(() {
        feed = result;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Friends Feed')),
      body: feed.isEmpty
          ? const Center(child: Text('No activity found.'))
          : ListView.builder(
              itemCount: feed.length,
              itemBuilder: (context, index) {
                final item = feed[index];
                return ListTile(
                  title: Text(item['description'] ?? 'No description'),
                  subtitle: Text('User ID: ${item['user_id']}'),
                );
              },
            ),
    );
  }
}
