import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  void _showFriendsDialog() {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Friends'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: const [
              TextField(
                decoration: InputDecoration(labelText: 'Search Friends'),
              ),
              SizedBox(height: 10),
              Text('Friend list loading...'), // Здесь можно список из API
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Close'),
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
          title: const Text('Friend Requests'),
          content: const Text('Здесь появятся входящие заявки в друзья'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Close'),
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
          icon: const Icon(Icons.menu),
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
            icon: const Icon(Icons.people_alt),
            onPressed: _showFriendsDialog,
          ),
          IconButton(
            icon: const Icon(Icons.notifications),
            onPressed: _showNotificationsDialog,
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
