import 'package:flutter/material.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final buttonStyle = ElevatedButton.styleFrom(
      backgroundColor: Colors.orange[700],
      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
    );

    return Scaffold(
      backgroundColor: Colors.orange[50],
      appBar: AppBar(
        backgroundColor: Colors.orange[700],
        title: const Text(
          'Healthy Summer',
          style: TextStyle(color: Colors.white),
        ),
      ),
      body: Center(
        child: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.local_florist, size: 80, color: Colors.orange[700]),
              const SizedBox(height: 20),
              Text(
                'Welcome to Healthy Summer!',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: Colors.orange[700],
                ),
              ),
              const SizedBox(height: 10),
              const Text(
                'Track your health, stay active, and enjoy summer!',
                textAlign: TextAlign.center,
                style: TextStyle(fontSize: 16),
              ),
              const SizedBox(height: 30),

              // Activities
              ElevatedButton.icon(
                style: buttonStyle,
                icon: const Icon(Icons.directions_run),
                label: const Text('Activities'),
                onPressed: () => Navigator.pushNamed(context, '/activities'),
              ),
              const SizedBox(height: 12),

              // Meals
              ElevatedButton.icon(
                style: buttonStyle,
                icon: const Icon(Icons.restaurant_menu),
                label: const Text('Meals'),
                onPressed: () => Navigator.pushNamed(context, '/meals'),
              ),
              const SizedBox(height: 12),

              // Water
              ElevatedButton.icon(
                style: buttonStyle,
                icon: const Icon(Icons.water_drop),
                label: const Text('Water Log'),
                onPressed: () => Navigator.pushNamed(context, '/water'),
              ),
              const SizedBox(height: 12),

              // Foods
              ElevatedButton.icon(
                style: buttonStyle,
                icon: const Icon(Icons.fastfood),
                label: const Text('Food Database'),
                onPressed: () => Navigator.pushNamed(context, '/foods'),
              ),
              const SizedBox(height: 12),

              // Stats placeholder
              ElevatedButton.icon(
                style: buttonStyle,
                icon: const Icon(Icons.bar_chart),
                label: const Text('Statistics (coming soon)'),
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Statistics feature in progress...'),
                    ),
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
