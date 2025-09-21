# VHS tape for gitsweeper demo
Output demo.gif

# Set theme and size  
Set FontSize 24
Set Width 1000
Set Height 600
Set Padding 20
Set WindowBar Colorful
Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10

# Show the title
Type "# gitsweeper - Clean up merged Git branches"
Enter
Sleep 2s

Type "# Let's see what gitsweeper can do"
Enter
Sleep 2s

# Show help
Type "./bin/gitsweeper --help"
Enter
Sleep 3s

# Show version
Type "./bin/gitsweeper --version"
Enter
Sleep 2s

# Demonstrate preview command
Type "echo '# Preview what branches would be cleaned:'"
Enter
Type "./bin/gitsweeper preview"
Enter
Sleep 3s

# Show options
Type "echo '# Options: --origin, --master, --skip, --force'"
Enter
Sleep 2s

# Show cleanup command
Type "echo '# To delete branches, use cleanup. It'll prompt us with a y/n prompt first unless we use -force"
Enter
Type "./bin/gitsweeper cleanup"
Enter
Sleep 2s
Ctrl+C

# Final message
Type "echo ''"
Enter
Type "echo '✨ Thats Gitsweeper: Fast, safe Git branch cleanup!'"
Enter
Sleep 3s