import 'package:flutter/material.dart';
import '../../services/api_service.dart';

class CreateChallengeScreen extends StatefulWidget {
  const CreateChallengeScreen({super.key});

  @override
  State<CreateChallengeScreen> createState() => _CreateChallengeScreenState();
}

class _CreateChallengeScreenState extends State<CreateChallengeScreen> {
  final ApiService api = ApiService();
  final TextEditingController _titleController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _goalValueController = TextEditingController();

  Future<void> submit() async {
    final title = _titleController.text;
    final desc = _descriptionController.text;
    final goal = int.tryParse(_goalValueController.text);

    if (title.isEmpty || desc.isEmpty || goal == null) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Fill in all fields')));
      return;
    }

    final success = await api.createChallenge(
      title: title,
      description: desc,
      goalValue: goal,
    );

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(success ? 'Challenge created!' : 'Error creating'),
      ),
    );

    if (success) {
      Navigator.pop(context);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Create Challenge')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(
              controller: _titleController,
              decoration: const InputDecoration(labelText: 'Title'),
            ),
            TextField(
              controller: _descriptionController,
              decoration: const InputDecoration(labelText: 'Description'),
            ),
            TextField(
              controller: _goalValueController,
              decoration: const InputDecoration(labelText: 'Goal Value'),
              keyboardType: TextInputType.number,
            ),
            const SizedBox(height: 20),
            ElevatedButton(onPressed: submit, child: const Text('Create')),
          ],
        ),
      ),
    );
  }
}
