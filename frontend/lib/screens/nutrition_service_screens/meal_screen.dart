import 'package:flutter/material.dart';
import '../../services/api_service.dart';
import './edit_meal_screen.dart';
import 'package:intl/intl.dart';

class MealsScreen extends StatefulWidget {
  const MealsScreen({Key? key}) : super(key: key);

  @override
  State<MealsScreen> createState() => _MealsScreenState();
}

class _MealsScreenState extends State<MealsScreen> {
  final ApiService _apiService = ApiService();
  List<dynamic> _meals = [];

  @override
  void initState() {
    super.initState();
    _loadMeals();
  }

  Future<void> _loadMeals() async {
    final meals = await _apiService.getMeals();
    if (meals != null) {
      setState(() {
        _meals = meals;
      });
    }
  }

  Future<void> _deleteMeal(int id) async {
    await _apiService.deleteMeal(id);
    _loadMeals();
  }

  String _formatDate(String iso) {
    try {
      return DateFormat('yyyy-MM-dd HH:mm').format(DateTime.parse(iso));
    } catch (_) {
      return iso;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        title: const Text('Meals'),
        backgroundColor: Colors.orange[700],
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(8),
        itemCount: _meals.length,
        itemBuilder: (context, index) {
          final meal = _meals[index];
          return Card(
            color: Colors.white,
            elevation: 2,
            margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 8),
            child: ListTile(
              title: Text(meal['description'] ?? 'No description'),
              subtitle: Text(
                'Calories: ${meal['calories']}\nTime: ${_formatDate(meal['meal_time'])}',
              ),
              isThreeLine: true,
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    icon: const Icon(Icons.edit),
                    onPressed: () => Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (_) => EditMealScreen(meal: meal),
                      ),
                    ).then((_) => _loadMeals()),
                  ),
                  IconButton(
                    icon: const Icon(Icons.delete),
                    onPressed: () => _deleteMeal(meal['id']),
                  ),
                ],
              ),
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () =>
            Navigator.pushNamed(context, '/add_meal').then((_) => _loadMeals()),
        backgroundColor: Colors.orange[700],
        child: const Icon(Icons.add),
      ),
    );
  }
}
