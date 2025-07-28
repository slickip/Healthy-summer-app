import 'package:flutter/material.dart';
import '../../services/api_service.dart';
import './edit_food_screen.dart';

class FoodsScreen extends StatefulWidget {
  const FoodsScreen({Key? key}) : super(key: key);

  @override
  State<FoodsScreen> createState() => _FoodsScreenState();
}

class _FoodsScreenState extends State<FoodsScreen> {
  final ApiService _apiService = ApiService();
  List<dynamic> _foods = [];

  String _val(dynamic value) {
    if (value == null) return '0';
    final parsed = double.tryParse(value.toString());
    return parsed != null ? parsed.toStringAsFixed(1) : '0';
  }

  @override
  void initState() {
    super.initState();
    _loadFoods();
  }

  Future<void> _loadFoods() async {
    final foods = await _apiService.getFoods();
    print('FoodsScreen - received foods: $foods');
    if (foods != null) {
      setState(() {
        _foods = foods;
      });
    }
  }

  Future<void> _deleteFood(int id) async {
    await _apiService.deleteFood(id);
    _loadFoods();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        title: const Text('Foods'),
        backgroundColor: Colors.orange[700],
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(8),
        itemCount: _foods.length,
        itemBuilder: (context, index) {
          final food = _foods[index];
          return Card(
            color: Colors.white,
            elevation: 2,
            margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 8),
            child: ListTile(
              title: Text(
                (food['name']?.toString().trim().isNotEmpty ?? false)
                    ? food['name']
                    : 'Unnamed',
              ),
              subtitle: Text(
                'Kcal: ${_val(food['callories_per_100g'])}, '
                'P: ${_val(food['proteins'])} '
                'F: ${_val(food['fats'])} '
                'C: ${_val(food['carbs'])}',
              ),

              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    icon: const Icon(Icons.edit),
                    onPressed: () => Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (_) => EditFoodScreen(food: food),
                      ),
                    ).then((_) => _loadFoods()),
                  ),
                  IconButton(
                    icon: const Icon(Icons.delete),
                    onPressed: () => _deleteFood(food['id']),
                  ),
                ],
              ),
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () =>
            Navigator.pushNamed(context, '/add_food').then((_) => _loadFoods()),
        backgroundColor: Colors.orange[700],
        child: const Icon(Icons.add),
      ),
    );
  }
}
