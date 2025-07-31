import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import '../services/api_service.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final ApiService api = ApiService();
  final storage = const FlutterSecureStorage();

  List<dynamic> searchResults = [];
  bool isSearching = false;
  final TextEditingController _searchController = TextEditingController();

  List<dynamic> friends = [];
  List<dynamic> requests = [];

  @override
  void initState() {
    super.initState();
    loadFriends();
    loadRequests();
  }

  Future<void> loadFriends() async {
    final result = await api.getFriendsList();
    if (result != null) {
      setState(() {
        friends = result;
      });
    }
  }

  Future<void> loadRequests() async {
    final result = await api.getIncomingRequests();
    if (result != null) {
      setState(() {
        requests = result;
      });
    }
  }

  Future<void> respondToRequest(int requestId, String action) async {
    final success = await api.respondToRequest(requestId, action);
    if (success) {
      loadRequests();
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('Request ${action}ed')));
    }
  }

  void _logout() async {
    await api.logout();
    if (!mounted) return;
    Navigator.pushReplacementNamed(context, '/login');
  }

  void _showFriendsDialog() {
    void searchUser() async {
      final query = _searchController.text;
      if (query.isEmpty) return;

      final result = await api.searchAllUsers(query);
      if (result != null) {
        setState(() {
          searchResults = result;
          isSearching = true;
        });
      }
    }

    void sendFriendRequest(int userId) async {
      final success = await api.sendFriendRequest(userId);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(success ? 'Request sent' : 'Failed to send')),
      );
    }

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          backgroundColor: Colors.orange[50],
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          title: const Text('Friends', style: TextStyle(color: Colors.orange)),
          content: SizedBox(
            width: double.maxFinite,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    labelText: 'Search user',
                    suffixIcon: IconButton(
                      icon: const Icon(Icons.search),
                      onPressed: searchUser,
                    ),
                  ),
                  onSubmitted: (_) => searchUser(),
                ),
                const SizedBox(height: 10),
                if (isSearching)
                  searchResults.isEmpty
                      ? const Text('This user was not found')
                      : SizedBox(
                          height: 300,
                          child: ListView.builder(
                            shrinkWrap: true,
                            itemCount: searchResults.length,
                            itemBuilder: (context, index) {
                              final user = searchResults[index];
                              return ListTile(
                                title: Text(
                                  user['display_name'] ?? 'User ${user['id']}',
                                ),
                                subtitle: Text(user['email'] ?? ''),
                                trailing: ElevatedButton(
                                  onPressed: () =>
                                      sendFriendRequest(user['id']),
                                  child: const Text('Add'),
                                ),
                              );
                            },
                          ),
                        )
                else
                  friends.isEmpty
                      ? const Text('No friends yet')
                      : ListView.builder(
                          shrinkWrap: true,
                          itemCount: friends.length,
                          itemBuilder: (context, index) {
                            final friend = friends[index];
                            return ListTile(
                              leading: const Icon(Icons.person),
                              title: Text(
                                friend['display_name'] ??
                                    'User ${friend['id']}',
                              ),
                              subtitle: Text(friend['email'] ?? ''),
                            );
                          },
                        ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                setState(() {
                  isSearching = false;
                  searchResults = [];
                });
                Navigator.pop(context);
              },
              child: const Text(
                'Close',
                style: TextStyle(color: Colors.orange),
              ),
            ),
          ],
        );
      },
    );
  }

  void _showNotificationsDialog() {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          backgroundColor: Colors.orange[50],
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          title: const Text(
            'Friend Requests',
            style: TextStyle(color: Colors.orange),
          ),
          content: SizedBox(
            width: double.maxFinite,
            child: requests.isEmpty
                ? const Text('No incoming requests')
                : ListView.builder(
                    shrinkWrap: true,
                    itemCount: requests.length,
                    itemBuilder: (context, index) {
                      final req = requests[index];
                      final sender = req['sender'];
                      final displayName = sender != null
                          ? sender['display_name'] ?? 'User ${req['sender_id']}'
                          : 'User ${req['sender_id']}';
                      final email = sender != null ? sender['email'] ?? '' : '';
                      return ListTile(
                        title: Text(displayName),
                        subtitle: Text(email),
                        trailing: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            IconButton(
                              icon: const Icon(
                                Icons.check,
                                color: Colors.green,
                              ),
                              onPressed: () =>
                                  respondToRequest(req['id'], 'accept'),
                            ),
                            IconButton(
                              icon: const Icon(Icons.close, color: Colors.red),
                              onPressed: () =>
                                  respondToRequest(req['id'], 'decline'),
                            ),
                          ],
                        ),
                      );
                    },
                  ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text(
                'Close',
                style: TextStyle(color: Colors.orange),
              ),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        leading: PopupMenuButton<String>(
          icon: const Icon(Icons.menu, color: Colors.white),
          onSelected: (value) {
            Navigator.pushNamed(context, value);
          },
          itemBuilder: (context) => [
            const PopupMenuItem(
              value: '/activities',
              child: Text('Activities'),
            ),
            const PopupMenuItem(value: '/meals', child: Text('Meals')),
            const PopupMenuItem(value: '/water', child: Text('Water Log')),
            const PopupMenuItem(value: '/foods', child: Text('Food Database')),
            const PopupMenuItem(
              value: '/challenge_list',
              child: Text('Challenges'),
            ),
          ],
        ),
        title: const Text(
          'Healthy Summer',
          style: TextStyle(color: Colors.white),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.people, color: Colors.white),
            onPressed: _showFriendsDialog,
          ),
          IconButton(
            icon: const Icon(Icons.notifications, color: Colors.white),
            onPressed: _showNotificationsDialog,
          ),
          IconButton(
            icon: const Icon(Icons.logout, color: Colors.white),
            onPressed: _logout,
          ),
        ],
      ),
      body: Column(
        children: [
          const SizedBox(height: 16),
          Text(
            'Friends Activity Feed',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: Colors.orange[800],
            ),
          ),
          const SizedBox(height: 8),
          Expanded(
            child: ListView(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: const [
                Card(
                  child: ListTile(title: Text('Friend A completed a workout')),
                ),
                Card(
                  child: ListTile(title: Text('Friend B drank 2L of water')),
                ),
                Card(
                  child: ListTile(title: Text('Friend C joined a challenge')),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
