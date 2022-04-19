# Habit tracker
 Is a simple command-line tool to keep track of your habits.

## Usage
To start a new habit, for example a daily piano habit, simply type `habit piano`. As you repeat your habit daily, 
retype `habit piano`. You can have multiple habits at the same time, simply type `habit surfing` to start a new surfing 
habit. You can list all your streaks with `habit`. Also, you can create a weekly habit by passing the `weekly` option
like so `habit -f weekly piano`
```
Usage of habit:
  -f string
    	Set the frecuency of the habit: daily(default), weekly. (shorthand) (default "daily")
  -frequency string
    	Set the frecuency of the habit: daily(default), weekly. (default "daily")
```

 ## Installation
Install by running `go install https://github.com/crmejia/habit/cmd/habit@latest`