import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class MessagesScreen extends StatefulWidget {
  const MessagesScreen({super.key});

  @override
  State<MessagesScreen> createState() => _MessagesScreenState();
}

class _MessagesScreenState extends State<MessagesScreen> {
  final ApiService api = ApiService();
  final TextEditingController _friendIdController = TextEditingController();
  final TextEditingController _messageController = TextEditingController();
  List<dynamic> messages = [];
  int? selectedFriendId;

  Future<void> loadMessages() async {
    if (selectedFriendId == null) return;
    final result = await api.getMessages(selectedFriendId!);
    if (result != null) {
      setState(() {
        messages = result;
      });
    }
  }

  Future<void> sendMessage() async {
    if (selectedFriendId == null || _messageController.text.isEmpty) return;
    final success = await api.sendMessage(
      selectedFriendId!,
      _messageController.text,
    );
    if (success) {
      _messageController.clear();
      await loadMessages();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Messages')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(
              controller: _friendIdController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(labelText: 'Friend ID'),
              onSubmitted: (value) {
                selectedFriendId = int.tryParse(value);
                loadMessages();
              },
            ),
            const SizedBox(height: 10),
            Expanded(
              child: ListView.builder(
                itemCount: messages.length,
                itemBuilder: (context, index) {
                  final msg = messages[index];
                  return ListTile(
                    title: Text(msg['content']),
                    subtitle: Text('From: ${msg['sender_id']}'),
                  );
                },
              ),
            ),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _messageController,
                    decoration: const InputDecoration(
                      hintText: 'Enter message',
                    ),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: sendMessage,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
