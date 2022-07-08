import java.util.*;

public class Main{

    public static void main(String[] args)
    {
        Scanner scanner = new Scanner(System.in);
        System.out.print("请输入你的名字:");
        String name = scanner.next();
        try {
            Thread.sleep(10);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        System.out.println("你好,"+name);
    }
}