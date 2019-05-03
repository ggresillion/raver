import {Injectable, NotFoundException} from '@nestjs/common';
import {Category} from './category.entity';
import {Repository} from 'typeorm';
import {InjectRepository} from '@nestjs/typeorm';

@Injectable()
export class CategoryService {

  constructor(
    @InjectRepository(Category)
    private readonly categoryRepository: Repository<Category>,
  ) {
  }

  public async findAllCategories(): Promise<Category[]> {
    return this.categoryRepository.find();
  }

  public async createCategory(category: Category) {
    return this.categoryRepository.insert(category);
  }

  public async deleteCategory(categoryId: number): Promise<Category> {
    const category = await this.categoryRepository.findOne(categoryId);
    if (!category) {
      throw new NotFoundException('category not found');
    }
    await this.categoryRepository.remove(category);
    return category;
  }
}
